"""LSTM Laura price prediction model.

Provides functions for downloading stock data, training/fine-tuning
an LSTM neural network, and predicting price movement
(drop / flat / rise) for a given ticker.
"""

import yfinance as yf
import numpy as np
import pandas as pd
import os
import time
import pickle
import sys
from sklearn.preprocessing import MinMaxScaler
from tensorflow.keras.models import Sequential, load_model
from tensorflow.keras.layers import LSTM, Dense, Dropout
from tensorflow.keras.utils import to_categorical
from loguru import logger

logger.remove()
logger.add(
    sys.stdout,
    format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
    level="DEBUG",
    colorize=True,
)

#настройки начального обучения
TICKER = 'AAPL'
MODEL_FILE = 'Laura.keras'
SCALER_FILE = 'scaler.pkl'
WINDOW = 10

def get_data(ticker):
    """Download one year of daily closing prices via yfinance.

    Computes percentage change and creates a 3-class target column:
    0 = drop (>1% down), 1 = flat, 2 = rise (>1% up).

    Args:
        ticker: Stock ticker symbol.

    Returns:
        DataFrame with 'Close', 'change', and 'target' columns.
    """
    df = yf.download(ticker, period='1y', progress=False)
    if df.empty:
        raise Exception("Нет данных")
    df = df[['Close']].copy()
    df['change'] = df['Close'].pct_change().shift(-1)
    df['target'] = np.where(df['change'] > 0.01, 2, np.where(df['change'] < -0.01, 0, 1))
    return df.dropna()

def sequences(df, scaler, fit=False):
    """Create sliding window sequences for LSTM training/prediction.

    Uses a window size of 10. When ``fit=True``, the scaler is fitted
    on the data; otherwise it transforms using an already-fitted scaler.

    Args:
        df: DataFrame with 'Close' and 'target' columns.
        scaler: sklearn MinMaxScaler instance.
        fit: Whether to fit the scaler on the data.

    Returns:
        Tuple of (X, y) where X is a 3D numpy array and y is one-hot encoded.
    """
    prices = df['Close'].values.reshape(-1, 1)
    
    if fit:
        data = scaler.fit_transform(prices)
    else:
        data = scaler.transform(prices)
        
    X, y = [], []
    targets = df['target'].values
    
    for i in range(WINDOW, len(data)):
        X.append(data[i-WINDOW:i])
        y.append(targets[i])
        
    return np.array(X), to_categorical(np.array(y), num_classes=3)

def create_model(shape):
    """Build and compile an LSTM model.

    Architecture: LSTM(64) -> Dropout(0.2) -> Dense(32, relu) -> Dense(3, softmax).
    Compiled with Adam optimizer and categorical crossentropy loss.

    Args:
        shape: Input shape tuple (window_size, n_features).

    Returns:
        Compiled Keras Sequential model.
    """
    model = Sequential([
        LSTM(64, input_shape=shape),
        Dropout(0.2),
        Dense(32, activation='relu'),
        Dense(3, activation='softmax')
    ])
    model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=['accuracy'])
    return model

def run(ticker):
    """Train or fine-tune the LSTM model for a given ticker.

    If a saved model (Laura.keras) and scaler (scaler.pkl) exist,
    fine-tunes for 2 epochs; otherwise trains from scratch for 15 epochs.

    Args:
        ticker: Stock ticker symbol.

    Returns:
        Tuple of (trained model, scaler, DataFrame).
    """
    logger.info(f"Загрузка {ticker}...")
    start = time.time()
    df = get_data(ticker)
    
    if os.path.exists(MODEL_FILE) and os.path.exists(SCALER_FILE):
        logger.info("Дообучение")
        model = load_model(MODEL_FILE)
        model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=['accuracy'])
        
        with open(SCALER_FILE, 'rb') as f:
            scaler = pickle.load(f)
        X, y = sequences(df, scaler, fit=False)
        model.fit(X, y, epochs=2, batch_size=32, verbose=0)
        
    else:
        logger.info("Модели нет")
        scaler = MinMaxScaler()
        X, y = sequences(df, scaler, fit=True)
        model = create_model((X.shape[1], X.shape[2]))
        model.fit(X, y, epochs=15, batch_size=32, verbose=1)
        with open(SCALER_FILE, 'wb') as f:
            pickle.dump(scaler, f)

    model.save(MODEL_FILE)
    logger.info(f"Готово за {time.time() - start:.2f}")
    
    return model, scaler, df

def predict(model, scaler, df):
    """Predict price movement direction for the next day.

    Uses the last WINDOW closing prices, scales them, and runs the model.

    Args:
        model: Trained Keras model.
        scaler: Fitted MinMaxScaler.
        df: DataFrame with historical price data.
    """
    last_prices = df['Close'].tail(WINDOW).values.reshape(-1, 1)
    scaled = scaler.transform(last_prices).reshape(1, WINDOW, 1)
    pred = model.predict(scaled, verbose=0)[0]
    
    labels = ["ПАДЕНИЕ", "ФЛЭТ", "РОСТ"]
    idx = np.argmax(pred)
    conf = pred[idx]

    logger.info(f"Направление: {labels[idx]}")
    logger.info(f"Уверенность: {conf:.1%}")
    # print(f"Вероятности: {pred[0]:.2f} / {pred[1]:.2f} / {pred[2]:.2f}")
    
    if conf > 0.5:
        logger.info(f"ИТОГ: {labels[idx]}")
    else:
        logger.info("хз")

if __name__ == "__main__":  
    while True:
        try:
            ticker_input = input()
            if ticker_input.lower() == 'exit':
                break
            m, s, d = run(ticker_input)
            predict(m, s, d)
            logger.info("-" * 20)
        except Exception as e:
            logger.error(f"Ошибка: {e}")
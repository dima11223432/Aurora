import json
from kafka import KafkaConsumer


#Дима вставь данные потом сюда сам, а то я хз какой у нас топик и сервак, на тестовом локалхосте у  сеня работало вс1ё норм
TOPIC = None
BOOTSTRAP_SERVERS = None
GROUP_ID = None
#TODO: Улучшить работу сервиса для вызова независимо от функции и доделать сам метод

def getMessage():
    consumer = KafkaConsumer(
        TOPIC,
        bootstrap_servers=BOOTSTRAP_SERVERS,
        group_id=GROUP_ID,
        auto_offset_reset='earliest',   
        enable_auto_commit=True,
        value_deserializer=lambda v: json.loads(v.decode('utf-8')), 
        key_deserializer=lambda k: k.decode('utf-8') if k else None
    )
    try:
        for message in consumer:
            data = message.value
            if all(key in data for key in ('id', 'date', 'text', 'channelUsername', 'channel_title', 'channelLink')):
                return data
            else:
                print("eror")
    except KeyboardInterrupt:
        pass
    finally:
        consumer.close()


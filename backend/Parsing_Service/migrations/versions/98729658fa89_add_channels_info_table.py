"""add_channels_info_table

Revision ID: 98729658fa89
Revises: 198c651c9d65
Create Date: 2026-04-25 13:54:26.222981

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "98729658fa89"
down_revision: Union[str, Sequence[str], None] = "198c651c9d65"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "channels_info",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column(
            "channel_id",
            sa.Integer,
            sa.ForeignKey("channels.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("category", sa.String(255), nullable=False),
        sa.UniqueConstraint("channel_id", "category", name="uq_channel_category"),
    )


def downgrade() -> None:
    op.drop_table("channels_info")

"""add_custom_user_channels

Revision ID: 0fe30db41e58
Revises: 599ed1df2acf
Create Date: 2026-04-27 19:34:49.170991

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "0fe30db41e58"
down_revision: Union[str, Sequence[str], None] = "599ed1df2acf"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "user_custom_parsing_channels",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column("user_id", sa.Integer, nullable=False),
        sa.Column("channel_username", sa.String(255), nullable=False),
    )


def downgrade() -> None:
    op.drop_table("user_custom_parsing_channels")

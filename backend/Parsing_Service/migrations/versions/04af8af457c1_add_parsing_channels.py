"""add_parsing_channels

Revision ID: 04af8af457c1
Revises: 70e338f5141d
Create Date: 2026-05-01 15:01:11.251146

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "04af8af457c1"
down_revision: Union[str, Sequence[str], None] = "70e338f5141d"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.drop_table("channels", if_exists=True)
    op.create_table(
        "channels",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column("username", sa.String(255), nullable=False, unique=True),
    )

    pass


def downgrade() -> None:
    op.drop_table("channels")
    pass

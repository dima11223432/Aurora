"""remane_channels_table_and_edit foreign key

Revision ID: 70e338f5141d
Revises: 6f3ea43bd222
Create Date: 2026-05-01 14:36:22.279630

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "70e338f5141d"
down_revision: Union[str, Sequence[str], None] = "6f3ea43bd222"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.drop_constraint(
        "channels_info_channel_id_fkey", "channels_info", type_="foreignkey"
    )

    op.create_foreign_key(
        "channels_info_channel_id_fkey",
        "channels_info",
        "default_channels",
        ["channel_id"],
        ["id"],
        ondelete="CASCADE",
    )


def downgrade() -> None:
    op.drop_constraint(
        "channels_info_channel_id_fkey", "channels_info", type_="foreignkey"
    )

    op.create_foreign_key(
        "channels_info_channel_id_fkey",
        "channels_info",
        "default_channels",
        ["channel_id"],
        ["id"],
        ondelete="CASCADE",
    )

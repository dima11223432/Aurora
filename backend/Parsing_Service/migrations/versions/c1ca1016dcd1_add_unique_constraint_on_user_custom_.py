"""add_unique_constraint_on_user_custom_channels

Revision ID: c1ca1016dcd1
Revises: 395b636b5a36
Create Date: 2026-05-01 22:25:09.586924

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "c1ca1016dcd1"
down_revision: Union[str, Sequence[str], None] = "395b636b5a36"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade():
    op.create_unique_constraint(
        "uq_user_channel",
        "user_custom_parsing_channels",
        ["user_id", "channel_username"],
    )


def downgrade():
    op.drop_constraint(
        "uq_user_channel", "user_custom_parsing_channels", type_="unique"
    )

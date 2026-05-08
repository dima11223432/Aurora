"""add_trigers_on_channels

Revision ID: 395b636b5a36
Revises: 04af8af457c1
Create Date: 2026-05-01 22:11:14.125587

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "395b636b5a36"
down_revision: Union[str, Sequence[str], None] = "04af8af457c1"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade():
    op.execute(
        """
        CREATE OR REPLACE FUNCTION notify_new_channel() RETURNS trigger AS $$
        BEGIN
            PERFORM pg_notify('new_channel_event', NEW.username);
            RETURN NEW;
        END;
        $$ LANGUAGE plpgsql;
    """
    )

    op.execute(
        """
        CREATE TRIGGER channel_insert_trigger
        AFTER INSERT ON channels
        FOR EACH ROW EXECUTE FUNCTION notify_new_channel();
    """
    )


def downgrade():
    op.execute("DROP TRIGGER IF EXISTS channel_insert_trigger ON channels;")
    op.execute("DROP FUNCTION IF EXISTS notify_new_channel();")

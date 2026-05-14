"""add_trigger on channels

Revision ID: a735170311f8
Revises: 0fe30db41e58
Create Date: 2026-04-27 20:15:10.847725

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "a735170311f8"
down_revision: Union[str, Sequence[str], None] = "0fe30db41e58"
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

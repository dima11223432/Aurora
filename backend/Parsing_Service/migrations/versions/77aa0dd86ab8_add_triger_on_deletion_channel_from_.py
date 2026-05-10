"""add_triger_on_deletion_channel_from_table

Revision ID: 77aa0dd86ab8
Revises: c1ca1016dcd1
Create Date: 2026-05-05 21:35:14.188200

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "77aa0dd86ab8"
down_revision: Union[str, Sequence[str], None] = "c1ca1016dcd1"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade():
    op.execute(
        """
        CREATE OR REPLACE FUNCTION notify_deleted_channel() RETURNS trigger AS $$
        BEGIN
            PERFORM pg_notify('deleted_channel_event', OLD.username);
            RETURN OLD;
        END;
        $$ LANGUAGE plpgsql;
    """
    )

    op.execute(
        """
        CREATE TRIGGER channel_deletion_trigger
        AFTER DELETE ON channels
        FOR EACH ROW EXECUTE FUNCTION notify_deleted_channel();
    """
    )


def downgrade():
    op.execute("DROP TRIGGER IF EXISTS channel_deletion_trigger ON channels;")
    op.execute("DROP FUNCTION IF EXISTS notify_deleted_channel();")

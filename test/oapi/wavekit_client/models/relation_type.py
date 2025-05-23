from enum import Enum


class RelationType(str, Enum):
    MATERIALIZEDVIEW = "materializedView"
    SINK = "sink"
    SOURCE = "source"
    SYSTEM_TABLE = "system table"
    TABLE = "table"

    def __str__(self) -> str:
        return str(self.value)

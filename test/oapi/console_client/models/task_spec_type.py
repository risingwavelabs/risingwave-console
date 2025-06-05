from enum import Enum


class TaskSpecType(str, Enum):
    AUTO_BACKUP = "auto-backup"
    AUTO_DIAGNOSTIC = "auto-diagnostic"
    DELETE_CLUSTER_DIAGNOSTIC = "delete-cluster-diagnostic"
    DELETE_OPAQUE_KEY = "delete-opaque-key"
    DELETE_SNAPSHOT = "delete-snapshot"

    def __str__(self) -> str:
        return str(self.value)

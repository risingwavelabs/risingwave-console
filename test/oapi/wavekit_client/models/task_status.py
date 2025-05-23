from enum import Enum


class TaskStatus(str, Enum):
    COMPLETED = "completed"
    FAILED = "failed"
    PAUSED = "paused"
    PENDING = "pending"

    def __str__(self) -> str:
        return str(self.value)

from enum import Enum


class EventSpecType(str, Enum):
    TASKCOMPLETED = "TaskCompleted"
    TASKERROR = "TaskError"

    def __str__(self) -> str:
        return str(self.value)

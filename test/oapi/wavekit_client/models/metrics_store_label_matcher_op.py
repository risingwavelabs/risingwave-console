from enum import Enum


class MetricsStoreLabelMatcherOp(str, Enum):
    EQ = "EQ"
    NEQ = "NEQ"
    NRE = "NRE"
    RE = "RE"

    def __str__(self) -> str:
        return str(self.value)

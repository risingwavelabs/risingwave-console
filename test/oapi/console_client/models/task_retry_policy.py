from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="TaskRetryPolicy")


@_attrs_define
class TaskRetryPolicy:
    """
    Attributes:
        interval (str): Interval of the retry policy, e.g. 1h, 1d, 1w, 1m
        always_retry_on_failure (Union[Unset, bool]): Whether to always retry the task on failure
    """

    interval: str
    always_retry_on_failure: Union[Unset, bool] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        interval = self.interval

        always_retry_on_failure = self.always_retry_on_failure

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "interval": interval,
            }
        )
        if always_retry_on_failure is not UNSET:
            field_dict["always_retry_on_failure"] = always_retry_on_failure

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        interval = d.pop("interval")

        always_retry_on_failure = d.pop("always_retry_on_failure", UNSET)

        task_retry_policy = cls(
            interval=interval,
            always_retry_on_failure=always_retry_on_failure,
        )

        task_retry_policy.additional_properties = d
        return task_retry_policy

    @property
    def additional_keys(self) -> list[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties

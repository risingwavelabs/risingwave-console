from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..models.event_spec_type import EventSpecType
from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.event_task_completed import EventTaskCompleted
    from ..models.event_task_error import EventTaskError


T = TypeVar("T", bound="EventSpec")


@_attrs_define
class EventSpec:
    """
    Attributes:
        type_ (EventSpecType):
        task_error (Union[Unset, EventTaskError]):
        task_completed (Union[Unset, EventTaskCompleted]):
    """

    type_: EventSpecType
    task_error: Union[Unset, "EventTaskError"] = UNSET
    task_completed: Union[Unset, "EventTaskCompleted"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        type_ = self.type_.value

        task_error: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.task_error, Unset):
            task_error = self.task_error.to_dict()

        task_completed: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.task_completed, Unset):
            task_completed = self.task_completed.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "type": type_,
            }
        )
        if task_error is not UNSET:
            field_dict["taskError"] = task_error
        if task_completed is not UNSET:
            field_dict["taskCompleted"] = task_completed

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.event_task_completed import EventTaskCompleted
        from ..models.event_task_error import EventTaskError

        d = dict(src_dict)
        type_ = EventSpecType(d.pop("type"))

        _task_error = d.pop("taskError", UNSET)
        task_error: Union[Unset, EventTaskError]
        if isinstance(_task_error, Unset):
            task_error = UNSET
        else:
            task_error = EventTaskError.from_dict(_task_error)

        _task_completed = d.pop("taskCompleted", UNSET)
        task_completed: Union[Unset, EventTaskCompleted]
        if isinstance(_task_completed, Unset):
            task_completed = UNSET
        else:
            task_completed = EventTaskCompleted.from_dict(_task_completed)

        event_spec = cls(
            type_=type_,
            task_error=task_error,
            task_completed=task_completed,
        )

        event_spec.additional_properties = d
        return event_spec

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

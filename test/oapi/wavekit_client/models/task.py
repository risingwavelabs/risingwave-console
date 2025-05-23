import datetime
from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

from ..models.task_status import TaskStatus
from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.task_attributes import TaskAttributes
    from ..models.task_spec import TaskSpec


T = TypeVar("T", bound="Task")


@_attrs_define
class Task:
    """
    Attributes:
        id (int):
        attributes (TaskAttributes):
        spec (TaskSpec):
        status (TaskStatus):
        created_at (datetime.datetime):
        updated_at (datetime.datetime):
        started_at (Union[Unset, datetime.datetime]):
    """

    id: int
    attributes: "TaskAttributes"
    spec: "TaskSpec"
    status: TaskStatus
    created_at: datetime.datetime
    updated_at: datetime.datetime
    started_at: Union[Unset, datetime.datetime] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        attributes = self.attributes.to_dict()

        spec = self.spec.to_dict()

        status = self.status.value

        created_at = self.created_at.isoformat()

        updated_at = self.updated_at.isoformat()

        started_at: Union[Unset, str] = UNSET
        if not isinstance(self.started_at, Unset):
            started_at = self.started_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "attributes": attributes,
                "spec": spec,
                "status": status,
                "createdAt": created_at,
                "updatedAt": updated_at,
            }
        )
        if started_at is not UNSET:
            field_dict["startedAt"] = started_at

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.task_attributes import TaskAttributes
        from ..models.task_spec import TaskSpec

        d = dict(src_dict)
        id = d.pop("ID")

        attributes = TaskAttributes.from_dict(d.pop("attributes"))

        spec = TaskSpec.from_dict(d.pop("spec"))

        status = TaskStatus(d.pop("status"))

        created_at = isoparse(d.pop("createdAt"))

        updated_at = isoparse(d.pop("updatedAt"))

        _started_at = d.pop("startedAt", UNSET)
        started_at: Union[Unset, datetime.datetime]
        if isinstance(_started_at, Unset):
            started_at = UNSET
        else:
            started_at = isoparse(_started_at)

        task = cls(
            id=id,
            attributes=attributes,
            spec=spec,
            status=status,
            created_at=created_at,
            updated_at=updated_at,
            started_at=started_at,
        )

        task.additional_properties = d
        return task

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

import datetime
from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

T = TypeVar("T", bound="DiagnosticData")


@_attrs_define
class DiagnosticData:
    """
    Attributes:
        id (int): Unique identifier of the diagnostic entry
        created_at (datetime.datetime): When the diagnostic data was collected
        content (str): Raw diagnostic data message containing system metrics and information
    """

    id: int
    created_at: datetime.datetime
    content: str
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        created_at = self.created_at.isoformat()

        content = self.content

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "createdAt": created_at,
                "content": content,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = d.pop("ID")

        created_at = isoparse(d.pop("createdAt"))

        content = d.pop("content")

        diagnostic_data = cls(
            id=id,
            created_at=created_at,
            content=content,
        )

        diagnostic_data.additional_properties = d
        return diagnostic_data

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

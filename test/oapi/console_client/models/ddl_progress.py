import datetime
from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

from ..types import UNSET, Unset

T = TypeVar("T", bound="DDLProgress")


@_attrs_define
class DDLProgress:
    """
    Attributes:
        id (int):
        statement (str):
        progress (str): Progress of the materialized view creation
        initialized_at (Union[Unset, datetime.datetime]): When the DDL operation was initialized
    """

    id: int
    statement: str
    progress: str
    initialized_at: Union[Unset, datetime.datetime] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        statement = self.statement

        progress = self.progress

        initialized_at: Union[Unset, str] = UNSET
        if not isinstance(self.initialized_at, Unset):
            initialized_at = self.initialized_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "statement": statement,
                "progress": progress,
            }
        )
        if initialized_at is not UNSET:
            field_dict["initializedAt"] = initialized_at

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = d.pop("ID")

        statement = d.pop("statement")

        progress = d.pop("progress")

        _initialized_at = d.pop("initializedAt", UNSET)
        initialized_at: Union[Unset, datetime.datetime]
        if isinstance(_initialized_at, Unset):
            initialized_at = UNSET
        else:
            initialized_at = isoparse(_initialized_at)

        ddl_progress = cls(
            id=id,
            statement=statement,
            progress=progress,
            initialized_at=initialized_at,
        )

        ddl_progress.additional_properties = d
        return ddl_progress

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

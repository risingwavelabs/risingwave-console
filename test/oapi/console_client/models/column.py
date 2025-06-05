from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

T = TypeVar("T", bound="Column")


@_attrs_define
class Column:
    """
    Attributes:
        name (str): Name of the column
        type_ (str): Data type of the column
        is_primary_key (bool): Whether the column is a primary key
        is_hidden (bool): Whether the column is hidden
    """

    name: str
    type_: str
    is_primary_key: bool
    is_hidden: bool
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        name = self.name

        type_ = self.type_

        is_primary_key = self.is_primary_key

        is_hidden = self.is_hidden

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "name": name,
                "type": type_,
                "isPrimaryKey": is_primary_key,
                "isHidden": is_hidden,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        name = d.pop("name")

        type_ = d.pop("type")

        is_primary_key = d.pop("isPrimaryKey")

        is_hidden = d.pop("isHidden")

        column = cls(
            name=name,
            type_=type_,
            is_primary_key=is_primary_key,
            is_hidden=is_hidden,
        )

        column.additional_properties = d
        return column

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

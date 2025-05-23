from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

if TYPE_CHECKING:
    from ..models.relation import Relation


T = TypeVar("T", bound="Schema")


@_attrs_define
class Schema:
    """
    Attributes:
        name (str): Name of the schema
        relations (list['Relation']):
    """

    name: str
    relations: list["Relation"]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        name = self.name

        relations = []
        for relations_item_data in self.relations:
            relations_item = relations_item_data.to_dict()
            relations.append(relations_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "name": name,
                "relations": relations,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.relation import Relation

        d = dict(src_dict)
        name = d.pop("name")

        relations = []
        _relations = d.pop("relations")
        for relations_item_data in _relations:
            relations_item = Relation.from_dict(relations_item_data)

            relations.append(relations_item)

        schema = cls(
            name=name,
            relations=relations,
        )

        schema.additional_properties = d
        return schema

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

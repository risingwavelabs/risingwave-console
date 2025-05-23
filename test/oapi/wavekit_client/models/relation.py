from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..models.relation_type import RelationType

if TYPE_CHECKING:
    from ..models.column import Column


T = TypeVar("T", bound="Relation")


@_attrs_define
class Relation:
    """
    Attributes:
        id (int): Unique identifier of the table
        name (str): Name of the table
        type_ (RelationType): Type of the relation
        columns (list['Column']): List of columns in the table
        dependencies (list[int]):
    """

    id: int
    name: str
    type_: RelationType
    columns: list["Column"]
    dependencies: list[int]
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        name = self.name

        type_ = self.type_.value

        columns = []
        for columns_item_data in self.columns:
            columns_item = columns_item_data.to_dict()
            columns.append(columns_item)

        dependencies = self.dependencies

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "name": name,
                "type": type_,
                "columns": columns,
                "dependencies": dependencies,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.column import Column

        d = dict(src_dict)
        id = d.pop("ID")

        name = d.pop("name")

        type_ = RelationType(d.pop("type"))

        columns = []
        _columns = d.pop("columns")
        for columns_item_data in _columns:
            columns_item = Column.from_dict(columns_item_data)

            columns.append(columns_item)

        dependencies = cast(list[int], d.pop("dependencies"))

        relation = cls(
            id=id,
            name=name,
            type_=type_,
            columns=columns,
            dependencies=dependencies,
        )

        relation.additional_properties = d
        return relation

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

from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.column import Column
    from ..models.query_response_rows_item import QueryResponseRowsItem


T = TypeVar("T", bound="QueryResponse")


@_attrs_define
class QueryResponse:
    """
    Attributes:
        columns (list['Column']):
        rows (list['QueryResponseRowsItem']):
        rows_affected (int): Number of rows affected by the query
        error (Union[Unset, str]): Error message if the query failed
    """

    columns: list["Column"]
    rows: list["QueryResponseRowsItem"]
    rows_affected: int
    error: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        columns = []
        for columns_item_data in self.columns:
            columns_item = columns_item_data.to_dict()
            columns.append(columns_item)

        rows = []
        for rows_item_data in self.rows:
            rows_item = rows_item_data.to_dict()
            rows.append(rows_item)

        rows_affected = self.rows_affected

        error = self.error

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "columns": columns,
                "rows": rows,
                "rowsAffected": rows_affected,
            }
        )
        if error is not UNSET:
            field_dict["error"] = error

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.column import Column
        from ..models.query_response_rows_item import QueryResponseRowsItem

        d = dict(src_dict)
        columns = []
        _columns = d.pop("columns")
        for columns_item_data in _columns:
            columns_item = Column.from_dict(columns_item_data)

            columns.append(columns_item)

        rows = []
        _rows = d.pop("rows")
        for rows_item_data in _rows:
            rows_item = QueryResponseRowsItem.from_dict(rows_item_data)

            rows.append(rows_item)

        rows_affected = d.pop("rowsAffected")

        error = d.pop("error", UNSET)

        query_response = cls(
            columns=columns,
            rows=rows,
            rows_affected=rows_affected,
            error=error,
        )

        query_response.additional_properties = d
        return query_response

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

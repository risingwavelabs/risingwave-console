from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

T = TypeVar("T", bound="TestClusterConnectionPayload")


@_attrs_define
class TestClusterConnectionPayload:
    """
    Attributes:
        host (str):
        sql_port (int):
        meta_port (int):
        http_port (int):
    """

    host: str
    sql_port: int
    meta_port: int
    http_port: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        host = self.host

        sql_port = self.sql_port

        meta_port = self.meta_port

        http_port = self.http_port

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "host": host,
                "sqlPort": sql_port,
                "metaPort": meta_port,
                "httpPort": http_port,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        host = d.pop("host")

        sql_port = d.pop("sqlPort")

        meta_port = d.pop("metaPort")

        http_port = d.pop("httpPort")

        test_cluster_connection_payload = cls(
            host=host,
            sql_port=sql_port,
            meta_port=meta_port,
            http_port=http_port,
        )

        test_cluster_connection_payload.additional_properties = d
        return test_cluster_connection_payload

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

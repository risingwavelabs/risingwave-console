from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="UpdateClusterRequest")


@_attrs_define
class UpdateClusterRequest:
    """
    Attributes:
        name (str):
        host (str):
        sql_port (int):
        meta_port (int):
        http_port (int):
        version (str):
        metrics_store_id (Union[Unset, int]): ID of the metrics store this cluster belongs to
    """

    name: str
    host: str
    sql_port: int
    meta_port: int
    http_port: int
    version: str
    metrics_store_id: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        name = self.name

        host = self.host

        sql_port = self.sql_port

        meta_port = self.meta_port

        http_port = self.http_port

        version = self.version

        metrics_store_id = self.metrics_store_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "name": name,
                "host": host,
                "sqlPort": sql_port,
                "metaPort": meta_port,
                "httpPort": http_port,
                "version": version,
            }
        )
        if metrics_store_id is not UNSET:
            field_dict["metricsStoreID"] = metrics_store_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        name = d.pop("name")

        host = d.pop("host")

        sql_port = d.pop("sqlPort")

        meta_port = d.pop("metaPort")

        http_port = d.pop("httpPort")

        version = d.pop("version")

        metrics_store_id = d.pop("metricsStoreID", UNSET)

        update_cluster_request = cls(
            name=name,
            host=host,
            sql_port=sql_port,
            meta_port=meta_port,
            http_port=http_port,
            version=version,
            metrics_store_id=metrics_store_id,
        )

        update_cluster_request.additional_properties = d
        return update_cluster_request

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

import datetime
from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

from ..types import UNSET, Unset

T = TypeVar("T", bound="Cluster")


@_attrs_define
class Cluster:
    """
    Attributes:
        id (int):
        org_id (int):
        name (str):
        host (str):
        sql_port (int):
        meta_port (int):
        http_port (int):
        version (str):
        created_at (datetime.datetime):
        updated_at (datetime.datetime):
        metrics_store_id (Union[Unset, int]): ID of the metrics store this cluster belongs to
    """

    id: int
    org_id: int
    name: str
    host: str
    sql_port: int
    meta_port: int
    http_port: int
    version: str
    created_at: datetime.datetime
    updated_at: datetime.datetime
    metrics_store_id: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        org_id = self.org_id

        name = self.name

        host = self.host

        sql_port = self.sql_port

        meta_port = self.meta_port

        http_port = self.http_port

        version = self.version

        created_at = self.created_at.isoformat()

        updated_at = self.updated_at.isoformat()

        metrics_store_id = self.metrics_store_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "OrgID": org_id,
                "name": name,
                "host": host,
                "sqlPort": sql_port,
                "metaPort": meta_port,
                "httpPort": http_port,
                "version": version,
                "createdAt": created_at,
                "updatedAt": updated_at,
            }
        )
        if metrics_store_id is not UNSET:
            field_dict["metricsStoreID"] = metrics_store_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = d.pop("ID")

        org_id = d.pop("OrgID")

        name = d.pop("name")

        host = d.pop("host")

        sql_port = d.pop("sqlPort")

        meta_port = d.pop("metaPort")

        http_port = d.pop("httpPort")

        version = d.pop("version")

        created_at = isoparse(d.pop("createdAt"))

        updated_at = isoparse(d.pop("updatedAt"))

        metrics_store_id = d.pop("metricsStoreID", UNSET)

        cluster = cls(
            id=id,
            org_id=org_id,
            name=name,
            host=host,
            sql_port=sql_port,
            meta_port=meta_port,
            http_port=http_port,
            version=version,
            created_at=created_at,
            updated_at=updated_at,
            metrics_store_id=metrics_store_id,
        )

        cluster.additional_properties = d
        return cluster

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

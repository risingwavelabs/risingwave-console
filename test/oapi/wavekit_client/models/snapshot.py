import datetime
from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

T = TypeVar("T", bound="Snapshot")


@_attrs_define
class Snapshot:
    """
    Attributes:
        id (int): Unique identifier of the snapshot
        cluster_id (int): ID of the cluster this snapshot belongs to
        name (str): Name of the snapshot
        created_at (datetime.datetime): Creation timestamp of the snapshot
    """

    id: int
    cluster_id: int
    name: str
    created_at: datetime.datetime
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        cluster_id = self.cluster_id

        name = self.name

        created_at = self.created_at.isoformat()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "ClusterID": cluster_id,
                "name": name,
                "createdAt": created_at,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = d.pop("ID")

        cluster_id = d.pop("ClusterID")

        name = d.pop("name")

        created_at = isoparse(d.pop("createdAt"))

        snapshot = cls(
            id=id,
            cluster_id=cluster_id,
            name=name,
            created_at=created_at,
        )

        snapshot.additional_properties = d
        return snapshot

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

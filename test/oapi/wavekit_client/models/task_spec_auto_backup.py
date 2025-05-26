from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

T = TypeVar("T", bound="TaskSpecAutoBackup")


@_attrs_define
class TaskSpecAutoBackup:
    """
    Attributes:
        cluster_id (int):
        retention_duration (str): Retention duration of the backup data, e.g. 1d, 1w, 1m, 1y
    """

    cluster_id: int
    retention_duration: str
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        cluster_id = self.cluster_id

        retention_duration = self.retention_duration

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "clusterID": cluster_id,
                "retentionDuration": retention_duration,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        cluster_id = d.pop("clusterID")

        retention_duration = d.pop("retentionDuration")

        task_spec_auto_backup = cls(
            cluster_id=cluster_id,
            retention_duration=retention_duration,
        )

        task_spec_auto_backup.additional_properties = d
        return task_spec_auto_backup

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

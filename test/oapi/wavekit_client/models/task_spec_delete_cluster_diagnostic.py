from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

T = TypeVar("T", bound="TaskSpecDeleteClusterDiagnostic")


@_attrs_define
class TaskSpecDeleteClusterDiagnostic:
    """
    Attributes:
        cluster_id (int):
        diagnostic_id (int):
    """

    cluster_id: int
    diagnostic_id: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        cluster_id = self.cluster_id

        diagnostic_id = self.diagnostic_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "clusterID": cluster_id,
                "diagnosticID": diagnostic_id,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        cluster_id = d.pop("clusterID")

        diagnostic_id = d.pop("diagnosticID")

        task_spec_delete_cluster_diagnostic = cls(
            cluster_id=cluster_id,
            diagnostic_id=diagnostic_id,
        )

        task_spec_delete_cluster_diagnostic.additional_properties = d
        return task_spec_delete_cluster_diagnostic

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

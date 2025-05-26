from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..models.task_spec_type import TaskSpecType
from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.task_spec_auto_backup import TaskSpecAutoBackup
    from ..models.task_spec_auto_diagnostic import TaskSpecAutoDiagnostic
    from ..models.task_spec_delete_cluster_diagnostic import TaskSpecDeleteClusterDiagnostic
    from ..models.task_spec_delete_opaque_key import TaskSpecDeleteOpaqueKey
    from ..models.task_spec_delete_snapshot import TaskSpecDeleteSnapshot


T = TypeVar("T", bound="TaskSpec")


@_attrs_define
class TaskSpec:
    """
    Attributes:
        type_ (TaskSpecType):
        auto_backup (Union[Unset, TaskSpecAutoBackup]):
        auto_diagnostic (Union[Unset, TaskSpecAutoDiagnostic]):
        delete_snapshot (Union[Unset, TaskSpecDeleteSnapshot]):
        delete_cluster_diagnostic (Union[Unset, TaskSpecDeleteClusterDiagnostic]):
        delete_opaque_key (Union[Unset, TaskSpecDeleteOpaqueKey]):
    """

    type_: TaskSpecType
    auto_backup: Union[Unset, "TaskSpecAutoBackup"] = UNSET
    auto_diagnostic: Union[Unset, "TaskSpecAutoDiagnostic"] = UNSET
    delete_snapshot: Union[Unset, "TaskSpecDeleteSnapshot"] = UNSET
    delete_cluster_diagnostic: Union[Unset, "TaskSpecDeleteClusterDiagnostic"] = UNSET
    delete_opaque_key: Union[Unset, "TaskSpecDeleteOpaqueKey"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        type_ = self.type_.value

        auto_backup: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.auto_backup, Unset):
            auto_backup = self.auto_backup.to_dict()

        auto_diagnostic: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.auto_diagnostic, Unset):
            auto_diagnostic = self.auto_diagnostic.to_dict()

        delete_snapshot: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.delete_snapshot, Unset):
            delete_snapshot = self.delete_snapshot.to_dict()

        delete_cluster_diagnostic: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.delete_cluster_diagnostic, Unset):
            delete_cluster_diagnostic = self.delete_cluster_diagnostic.to_dict()

        delete_opaque_key: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.delete_opaque_key, Unset):
            delete_opaque_key = self.delete_opaque_key.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "type": type_,
            }
        )
        if auto_backup is not UNSET:
            field_dict["autoBackup"] = auto_backup
        if auto_diagnostic is not UNSET:
            field_dict["autoDiagnostic"] = auto_diagnostic
        if delete_snapshot is not UNSET:
            field_dict["deleteSnapshot"] = delete_snapshot
        if delete_cluster_diagnostic is not UNSET:
            field_dict["deleteClusterDiagnostic"] = delete_cluster_diagnostic
        if delete_opaque_key is not UNSET:
            field_dict["deleteOpaqueKey"] = delete_opaque_key

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.task_spec_auto_backup import TaskSpecAutoBackup
        from ..models.task_spec_auto_diagnostic import TaskSpecAutoDiagnostic
        from ..models.task_spec_delete_cluster_diagnostic import TaskSpecDeleteClusterDiagnostic
        from ..models.task_spec_delete_opaque_key import TaskSpecDeleteOpaqueKey
        from ..models.task_spec_delete_snapshot import TaskSpecDeleteSnapshot

        d = dict(src_dict)
        type_ = TaskSpecType(d.pop("type"))

        _auto_backup = d.pop("autoBackup", UNSET)
        auto_backup: Union[Unset, TaskSpecAutoBackup]
        if isinstance(_auto_backup, Unset):
            auto_backup = UNSET
        else:
            auto_backup = TaskSpecAutoBackup.from_dict(_auto_backup)

        _auto_diagnostic = d.pop("autoDiagnostic", UNSET)
        auto_diagnostic: Union[Unset, TaskSpecAutoDiagnostic]
        if isinstance(_auto_diagnostic, Unset):
            auto_diagnostic = UNSET
        else:
            auto_diagnostic = TaskSpecAutoDiagnostic.from_dict(_auto_diagnostic)

        _delete_snapshot = d.pop("deleteSnapshot", UNSET)
        delete_snapshot: Union[Unset, TaskSpecDeleteSnapshot]
        if isinstance(_delete_snapshot, Unset):
            delete_snapshot = UNSET
        else:
            delete_snapshot = TaskSpecDeleteSnapshot.from_dict(_delete_snapshot)

        _delete_cluster_diagnostic = d.pop("deleteClusterDiagnostic", UNSET)
        delete_cluster_diagnostic: Union[Unset, TaskSpecDeleteClusterDiagnostic]
        if isinstance(_delete_cluster_diagnostic, Unset):
            delete_cluster_diagnostic = UNSET
        else:
            delete_cluster_diagnostic = TaskSpecDeleteClusterDiagnostic.from_dict(_delete_cluster_diagnostic)

        _delete_opaque_key = d.pop("deleteOpaqueKey", UNSET)
        delete_opaque_key: Union[Unset, TaskSpecDeleteOpaqueKey]
        if isinstance(_delete_opaque_key, Unset):
            delete_opaque_key = UNSET
        else:
            delete_opaque_key = TaskSpecDeleteOpaqueKey.from_dict(_delete_opaque_key)

        task_spec = cls(
            type_=type_,
            auto_backup=auto_backup,
            auto_diagnostic=auto_diagnostic,
            delete_snapshot=delete_snapshot,
            delete_cluster_diagnostic=delete_cluster_diagnostic,
            delete_opaque_key=delete_opaque_key,
        )

        task_spec.additional_properties = d
        return task_spec

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

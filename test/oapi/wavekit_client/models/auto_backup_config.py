from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

T = TypeVar("T", bound="AutoBackupConfig")


@_attrs_define
class AutoBackupConfig:
    """
    Attributes:
        enabled (bool): Whether automatic snapshots are enabled
        cron_expression (str): Cron expression for automatic snapshots (e.g., '0 0 * * *')
        retention_duration (str): How long to retain automatic snapshots (e.g., '1d', '7d', '14d', '30d', '90d')
    """

    enabled: bool
    cron_expression: str
    retention_duration: str
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        enabled = self.enabled

        cron_expression = self.cron_expression

        retention_duration = self.retention_duration

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "enabled": enabled,
                "cronExpression": cron_expression,
                "retentionDuration": retention_duration,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        enabled = d.pop("enabled")

        cron_expression = d.pop("cronExpression")

        retention_duration = d.pop("retentionDuration")

        auto_backup_config = cls(
            enabled=enabled,
            cron_expression=cron_expression,
            retention_duration=retention_duration,
        )

        auto_backup_config.additional_properties = d
        return auto_backup_config

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

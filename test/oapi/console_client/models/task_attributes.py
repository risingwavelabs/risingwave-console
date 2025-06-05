from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.task_cronjob import TaskCronjob
    from ..models.task_retry_policy import TaskRetryPolicy


T = TypeVar("T", bound="TaskAttributes")


@_attrs_define
class TaskAttributes:
    """
    Attributes:
        org_id (Union[Unset, int]): If the task is created by a user, this field will be the organization ID of the user
        timeout (Union[Unset, str]): Timeout of the task, e.g. 1h, 1d, 1w, 1m
        cronjob (Union[Unset, TaskCronjob]):
        retry_policy (Union[Unset, TaskRetryPolicy]):
    """

    org_id: Union[Unset, int] = UNSET
    timeout: Union[Unset, str] = UNSET
    cronjob: Union[Unset, "TaskCronjob"] = UNSET
    retry_policy: Union[Unset, "TaskRetryPolicy"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        org_id = self.org_id

        timeout = self.timeout

        cronjob: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.cronjob, Unset):
            cronjob = self.cronjob.to_dict()

        retry_policy: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.retry_policy, Unset):
            retry_policy = self.retry_policy.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if org_id is not UNSET:
            field_dict["orgID"] = org_id
        if timeout is not UNSET:
            field_dict["timeout"] = timeout
        if cronjob is not UNSET:
            field_dict["cronjob"] = cronjob
        if retry_policy is not UNSET:
            field_dict["retryPolicy"] = retry_policy

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.task_cronjob import TaskCronjob
        from ..models.task_retry_policy import TaskRetryPolicy

        d = dict(src_dict)
        org_id = d.pop("orgID", UNSET)

        timeout = d.pop("timeout", UNSET)

        _cronjob = d.pop("cronjob", UNSET)
        cronjob: Union[Unset, TaskCronjob]
        if isinstance(_cronjob, Unset):
            cronjob = UNSET
        else:
            cronjob = TaskCronjob.from_dict(_cronjob)

        _retry_policy = d.pop("retryPolicy", UNSET)
        retry_policy: Union[Unset, TaskRetryPolicy]
        if isinstance(_retry_policy, Unset):
            retry_policy = UNSET
        else:
            retry_policy = TaskRetryPolicy.from_dict(_retry_policy)

        task_attributes = cls(
            org_id=org_id,
            timeout=timeout,
            cronjob=cronjob,
            retry_policy=retry_policy,
        )

        task_attributes.additional_properties = d
        return task_attributes

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

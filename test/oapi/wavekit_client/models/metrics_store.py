import datetime
from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.metrics_store_label_matcher import MetricsStoreLabelMatcher
    from ..models.metrics_store_spec import MetricsStoreSpec


T = TypeVar("T", bound="MetricsStore")


@_attrs_define
class MetricsStore:
    """
    Attributes:
        id (int):
        name (str):
        created_at (datetime.datetime):
        spec (Union[Unset, MetricsStoreSpec]):
        default_labels (Union[Unset, list['MetricsStoreLabelMatcher']]):
    """

    id: int
    name: str
    created_at: datetime.datetime
    spec: Union[Unset, "MetricsStoreSpec"] = UNSET
    default_labels: Union[Unset, list["MetricsStoreLabelMatcher"]] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        name = self.name

        created_at = self.created_at.isoformat()

        spec: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.spec, Unset):
            spec = self.spec.to_dict()

        default_labels: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.default_labels, Unset):
            default_labels = []
            for componentsschemas_metrics_store_label_matcher_list_item_data in self.default_labels:
                componentsschemas_metrics_store_label_matcher_list_item = (
                    componentsschemas_metrics_store_label_matcher_list_item_data.to_dict()
                )
                default_labels.append(componentsschemas_metrics_store_label_matcher_list_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "name": name,
                "createdAt": created_at,
            }
        )
        if spec is not UNSET:
            field_dict["spec"] = spec
        if default_labels is not UNSET:
            field_dict["defaultLabels"] = default_labels

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.metrics_store_label_matcher import MetricsStoreLabelMatcher
        from ..models.metrics_store_spec import MetricsStoreSpec

        d = dict(src_dict)
        id = d.pop("ID")

        name = d.pop("name")

        created_at = isoparse(d.pop("createdAt"))

        _spec = d.pop("spec", UNSET)
        spec: Union[Unset, MetricsStoreSpec]
        if isinstance(_spec, Unset):
            spec = UNSET
        else:
            spec = MetricsStoreSpec.from_dict(_spec)

        default_labels = []
        _default_labels = d.pop("defaultLabels", UNSET)
        for componentsschemas_metrics_store_label_matcher_list_item_data in _default_labels or []:
            componentsschemas_metrics_store_label_matcher_list_item = MetricsStoreLabelMatcher.from_dict(
                componentsschemas_metrics_store_label_matcher_list_item_data
            )

            default_labels.append(componentsschemas_metrics_store_label_matcher_list_item)

        metrics_store = cls(
            id=id,
            name=name,
            created_at=created_at,
            spec=spec,
            default_labels=default_labels,
        )

        metrics_store.additional_properties = d
        return metrics_store

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

from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="TestDatabaseConnectionPayload")


@_attrs_define
class TestDatabaseConnectionPayload:
    """
    Attributes:
        cluster_id (int):
        username (str):
        database (str):
        password (Union[Unset, str]):
    """

    cluster_id: int
    username: str
    database: str
    password: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        cluster_id = self.cluster_id

        username = self.username

        database = self.database

        password = self.password

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "clusterID": cluster_id,
                "username": username,
                "database": database,
            }
        )
        if password is not UNSET:
            field_dict["password"] = password

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        cluster_id = d.pop("clusterID")

        username = d.pop("username")

        database = d.pop("database")

        password = d.pop("password", UNSET)

        test_database_connection_payload = cls(
            cluster_id=cluster_id,
            username=username,
            database=database,
            password=password,
        )

        test_database_connection_payload.additional_properties = d
        return test_database_connection_payload

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

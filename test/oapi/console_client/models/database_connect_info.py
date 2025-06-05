from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="DatabaseConnectInfo")


@_attrs_define
class DatabaseConnectInfo:
    """
    Attributes:
        name (str): Name of the database
        cluster_id (int): ID of the cluster this database belongs to
        username (str): Database username
        database (str): Database name
        password (Union[Unset, str]): Database password (optional)
    """

    name: str
    cluster_id: int
    username: str
    database: str
    password: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        name = self.name

        cluster_id = self.cluster_id

        username = self.username

        database = self.database

        password = self.password

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "name": name,
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
        name = d.pop("name")

        cluster_id = d.pop("clusterID")

        username = d.pop("username")

        database = d.pop("database")

        password = d.pop("password", UNSET)

        database_connect_info = cls(
            name=name,
            cluster_id=cluster_id,
            username=username,
            database=database,
            password=password,
        )

        database_connect_info.additional_properties = d
        return database_connect_info

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

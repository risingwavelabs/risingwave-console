import datetime
from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field
from dateutil.parser import isoparse

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.schema import Schema


T = TypeVar("T", bound="Database")


@_attrs_define
class Database:
    """
    Attributes:
        id (int): Unique identifier of the database
        name (str): Name of the database
        cluster_id (int): ID of the cluster this database belongs to
        org_id (int): ID of the organization this database belongs to
        username (str): Database username
        database (str): Database name
        created_at (datetime.datetime): Creation timestamp
        updated_at (datetime.datetime): Last update timestamp
        password (Union[Unset, str]): Database password (optional)
        schemas (Union[Unset, list['Schema']]): List of schemas in the database
    """

    id: int
    name: str
    cluster_id: int
    org_id: int
    username: str
    database: str
    created_at: datetime.datetime
    updated_at: datetime.datetime
    password: Union[Unset, str] = UNSET
    schemas: Union[Unset, list["Schema"]] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        name = self.name

        cluster_id = self.cluster_id

        org_id = self.org_id

        username = self.username

        database = self.database

        created_at = self.created_at.isoformat()

        updated_at = self.updated_at.isoformat()

        password = self.password

        schemas: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.schemas, Unset):
            schemas = []
            for schemas_item_data in self.schemas:
                schemas_item = schemas_item_data.to_dict()
                schemas.append(schemas_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "ID": id,
                "name": name,
                "clusterID": cluster_id,
                "OrgID": org_id,
                "username": username,
                "database": database,
                "createdAt": created_at,
                "updatedAt": updated_at,
            }
        )
        if password is not UNSET:
            field_dict["password"] = password
        if schemas is not UNSET:
            field_dict["schemas"] = schemas

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.schema import Schema

        d = dict(src_dict)
        id = d.pop("ID")

        name = d.pop("name")

        cluster_id = d.pop("clusterID")

        org_id = d.pop("OrgID")

        username = d.pop("username")

        database = d.pop("database")

        created_at = isoparse(d.pop("createdAt"))

        updated_at = isoparse(d.pop("updatedAt"))

        password = d.pop("password", UNSET)

        schemas = []
        _schemas = d.pop("schemas", UNSET)
        for schemas_item_data in _schemas or []:
            schemas_item = Schema.from_dict(schemas_item_data)

            schemas.append(schemas_item)

        database = cls(
            id=id,
            name=name,
            cluster_id=cluster_id,
            org_id=org_id,
            username=username,
            database=database,
            created_at=created_at,
            updated_at=updated_at,
            password=password,
            schemas=schemas,
        )

        database.additional_properties = d
        return database

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

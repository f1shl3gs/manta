import React from "react";
import { ResourceCard } from "@influxdata/clockface";
import { useHistory } from "react-router-dom";
import { useOrgID } from "../../shared/state/organization/organization";

interface Props {
  id: string
  name: string
  desc: string
  updatedAt: string
}

const DashboardCard: React.FC<Props> = props => {
  const {
    id,
    name,
    desc,
    updatedAt
  } = props;
  const history = useHistory();
  const orgID = useOrgID();

  const contextMenu = (): JSX.Element => {
    return (
      <div>Context Menu</div>
    );
  };

  return (
    <ResourceCard
      key={`dashboard-id--${id}`}
      contextMenu={contextMenu()}
    >
      <ResourceCard.EditableName
        onUpdate={(v) => console.log("update dashboard name", v)}
        onClick={() => history.push(`/orgs/${orgID}/dashboards/${id}`)}
        name={name}
      />

      <ResourceCard.EditableDescription
        onUpdate={(desc) => console.log("update desc", desc)}
        description={desc}
        placeholder={`Describe ${name}`}
      />
      <ResourceCard.Meta>
        {`Last Modified: ${updatedAt}`}
      </ResourceCard.Meta>
    </ResourceCard>
  );
};

export default DashboardCard;

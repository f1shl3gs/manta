import React from "react";
import { ResourceCard } from "@influxdata/clockface";

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
        onUpdate={() => console.log("update dashboard name")}
        onClick={() => console.log("on click", name)}
        name={name}
      />

      <ResourceCard.EditableDescription
        onUpdate={() => console.log("update desc")}
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

import React from "react";
import { Dashboard } from "types";
import DashboardCard from "./DashboardCard";

// Styles
import './DashboardCardsGrid.scss';

interface Props {
  dashboards: Dashboard[]
}

const DashboardCards: React.FC<Props> = props => {
  const {
    dashboards
  } = props;

  return (
    <div style={{height: '100%', display: 'grid'}}>
      <div className={"dashboards-card-grid"}>
        {
          dashboards.map(d => (
            <DashboardCard
              key={d.id}
              id={d.id}
              name={d.name}
              desc={d.desc}
              updatedAt={d.updated}
            />
          ))
        }
      </div>
    </div>
  );
};

export default DashboardCards;

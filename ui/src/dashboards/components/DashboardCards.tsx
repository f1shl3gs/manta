import React from "react";
import DashboardCard from "./DashboardCard";

// Styles
import './DashboardCardsGrid.scss';
import { useDashboards } from "../state";

const DashboardCards: React.FC = () => {
  const {dashboards} = useDashboards()

  return (
    <div style={{height: '100%', display: 'grid'}}>
      <div className={"dashboards-card-grid"}>
        {
          dashboards?.map(d => (
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

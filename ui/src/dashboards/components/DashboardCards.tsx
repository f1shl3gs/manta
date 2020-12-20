import React, { useCallback } from "react";
import DashboardCard from "./DashboardCard";

// Styles
import "./DashboardCardsGrid.scss";
import { useDashboards } from "../state";
import { useFetch } from "use-http";

const DashboardCards: React.FC = () => {
  const { dashboards, refresh } = useDashboards();

  const { del } = useFetch(`/api/v1/dashboards`, {});
  const onDeleteDashboard = useCallback((id: string) => {
    del(id)
      .then(() => {
        refresh();
      })
      .catch(err => {
        console.log("delete dashboard err", err);
      });
  }, [del]);

  return (
    <div style={{ height: "100%", display: "grid" }}>
      <div className={"dashboards-card-grid"}>
        {
          dashboards?.map(d => (
            <DashboardCard
              key={d.id}
              id={d.id}
              name={d.name}
              desc={d.desc}
              updatedAt={d.updated}
              onDeleteDashboard={onDeleteDashboard}
            />
          ))
        }
      </div>
    </div>
  );
};

export default DashboardCards;

// Libraries
import React, { useCallback } from 'react';
import moment from 'moment';

// Components
import {
  Button,
  ComponentColor,
  ComponentSize,
  EmptyState,
  IconFont,
  ResourceCard,
  ResourceList,
  SpinnerContainer,
  TechnoSpinner
} from '@influxdata/clockface';

// Hooks
import { useOrgID } from 'shared/hooks/useOrg';
import { useOtcls } from './state';
import { useFetch } from 'use-http';
import { useHistory } from 'react-router-dom';

// Types
import { Otcl } from 'types/otcl';

const Otcls: React.FC = () => {
  const orgID = useOrgID();
  const history = useHistory();
  const { otcls, rds, reload } = useOtcls();
  const { del } = useFetch(`/api/v1/otcls`, {});

  const context = useCallback((id: string): JSX.Element => {
    return (
      <Button
        icon={IconFont.Trash}
        text="Delete"
        size={ComponentSize.ExtraSmall}
        color={ComponentColor.Danger}
        onClick={() => {
          del(`${id}`)
            .then(() => {
              console.log('delete', id, 'success');
              reload();
            })
            .catch(() => {
              console.log('failed');
            });
        }}
      />
    );
  }, []);

  return (
    <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
      <ResourceList>
        <ResourceList.Body
          emptyState={
            <EmptyState size={ComponentSize.Large}>
              <EmptyState.Text>
                Looks like this Org doesn't have any <b>Otcl</b> configs
              </EmptyState.Text>
            </EmptyState>
          }
        >
          {otcls?.map((item: Otcl) => (
            <ResourceCard key={item.id} contextMenu={context(item.id)}>
              <ResourceCard.Name
                name={item.name}
                onClick={() => {
                  history.push(`/orgs/${orgID}/otcls/${item.id}`);
                }}
              />
              <ResourceCard.Meta>
                <span>updated: {moment(item.updated).fromNow()}</span>
                <span>desc: {item.desc}</span>
              </ResourceCard.Meta>
            </ResourceCard>
          ))}
        </ResourceList.Body>
      </ResourceList>
    </SpinnerContainer>
  );
};

export default Otcls;

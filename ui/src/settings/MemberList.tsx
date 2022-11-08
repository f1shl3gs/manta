// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {useResources} from 'src/shared/components/GetResources'
import {
  ComponentColor,
  IconFont,
  ResourceCard,
  ResourceList,
  Sort,
} from '@influxdata/clockface'
import {fromNow} from 'src/shared/duration'
import Context from 'src/shared/components/context_menu/Context'

const MemberList: FunctionComponent = () => {
  const {resources: members} = useResources()

  return (
    <ResourceList>
      <ResourceList.Header>
        <ResourceList.Sorter name="Name" sortKey="name" sort={Sort.Ascending} />
      </ResourceList.Header>

      <ResourceList.Body emptyState={<p>empty</p>}>
        {members.map(member => (
          <ResourceCard
            key={member.id}
            contextMenu={
              <Context>
                <Context.Menu
                  icon={IconFont.Trash_New}
                  color={ComponentColor.Danger}
                >
                  <Context.Item
                    label="Delete"
                    action={() => console.log('delete')}
                    testID="confirmation-button"
                  />
                </Context.Menu>
              </Context>
            }
          >
            <ResourceCard.EditableName
              onUpdate={() => console.log('name')}
              name={member.name}
              noNameString={''}
              buttonTestID="editable-name"
              inputTestID="input-field"
            />

            <ResourceCard.Meta>
              <>Last updated {fromNow(member.updated)}</>
            </ResourceCard.Meta>
          </ResourceCard>
        ))}
      </ResourceList.Body>
    </ResourceList>
  )
}

export default MemberList

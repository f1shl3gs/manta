// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ComponentColor,
  IconFont,
  ResourceCard,
  ResourceList,
  Sort,
} from '@influxdata/clockface'
import {fromNow} from 'src/shared/utils/duration'
import Context from 'src/shared/components/context_menu/Context'
import {useSelector} from 'react-redux'
import {User} from 'src/types/user'
import {getAll} from 'src/resources/selectors'
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'

const MemberList: FunctionComponent = () => {
  const members = useSelector((state: AppState) =>
    getAll<User>(state, ResourceType.Members)
  )

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

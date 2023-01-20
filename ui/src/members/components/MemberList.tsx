// Libraries
import React, {FunctionComponent} from 'react'
import {useSelector} from 'react-redux'

// Components
import {ResourceList, Sort} from '@influxdata/clockface'
import MemberCard from 'src/members/components/MemberCard'

// Types
import {User} from 'src/types/user'
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'

// Selectors
import {getAll} from 'src/resources/selectors'

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
          <MemberCard key={member.id} member={member} />
        ))}
      </ResourceList.Body>
    </ResourceList>
  )
}

export default MemberList

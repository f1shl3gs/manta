// Libraries
import React, {FunctionComponent} from 'react'

// Components
import MemberList from 'src/members/components/MemberList'
import GetResources from 'src/resources/components/GetResources'

// Types
import {ResourceType} from 'src/types/resources'

const MembersIndex: FunctionComponent = () => {
  return (
    <GetResources resources={[ResourceType.Members]}>
      <MemberList />
    </GetResources>
  )
}

export default MembersIndex

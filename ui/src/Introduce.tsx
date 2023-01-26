// Libraries
import React, {FunctionComponent} from 'react'
import {useSelector} from 'react-redux'

// Components
import {Page} from '@influxdata/clockface'

// Selectors
import {getOrg} from './organizations/selectors'

const Introduce: FunctionComponent = () => {
  const org = useSelector(getOrg)

  return (
    <Page titleTag={'Introduce'}>
      <Page.Header fullWidth={true}>
        <Page.Title title={'Getting started'} />
      </Page.Header>

      <Page.Contents testID={'introduction--page'}>{org.name}</Page.Contents>
    </Page>
  )
}

export default Introduce

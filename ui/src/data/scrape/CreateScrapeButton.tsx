import React, {FunctionComponent} from 'react'
import useFetch from 'src/shared/useFetch'
import {Button, ComponentColor} from '@influxdata/clockface'
import {useOrganization} from 'src/organizations/useOrganizations'

const CreateScrapeButton: FunctionComponent = () => {
  const {id: orgID} = useOrganization()
  const {run: create} = useFetch(`/api/v1/scrapes`, {
    method: 'POST',
  })

  return (
    <Button
      text={'Create Scrape'}
      color={ComponentColor.Primary}
      onClick={() =>
        create({
          orgID,
          name: 'selfstat',
          labels: {
            foo: 'bar',
          },
          targets: ['localhost:8088'],
        })
      }
    />
  )
}

export default CreateScrapeButton

// Libraries
import React, {useState} from 'react'

// Components
import {Columns, Grid} from '@influxdata/clockface'
import PackageSearch from './components/PackageSearch'
import FlameGraphControlBar from './FlameGraphControlBar'

const ProfilePanelHeader: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('')

  return (
    <Grid>
      <Grid.Column
        widthLG={Columns.Three}
        widthXS={Columns.Three}
        widthSM={Columns.Three}
        widthMD={Columns.Three}
      >
        <PackageSearch searchTerm={searchTerm} setSearchTerm={setSearchTerm} />
      </Grid.Column>

      <Grid.Column
        offsetLG={Columns.Four}
        offsetMD={Columns.Four}
        offsetSM={Columns.Four}
        offsetXS={Columns.Four}
        widthLG={Columns.Three}
        widthMD={Columns.Three}
        widthSM={Columns.Three}
        widthXS={Columns.Three}
      >
        <FlameGraphControlBar />
      </Grid.Column>
    </Grid>
  )
}

export default ProfilePanelHeader

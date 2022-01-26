import { Link } from 'react-router-dom';
import { useQuery, gql } from '@apollo/client';
import Metadata from './Metadata.js';
import CreateStation from './CreateStation.js';

function StationList() {
  const { loading, error, data } = useQuery(gql`
    query StationList {
      allStations {
	edges {
	  node {
	    id
            slug
            trackEventsByStationId(condition:{ action: "played"}, orderBy: CREATED_AT_DESC, first: 1) {
              nodes {
        	id
                trackByTrackId {
                  id
                  filename
                  metadata
                }
              }
            }
	  }
	}
      }
    }`, {
      pollInterval: 10000,
  });

  if (loading) {
    return "spinner"
  }

  return (
    <div>
      <h2>Stations</h2>
      {data.allStations.edges.map(({ node }) => (
        <div key={node.id}>
          <Link to={`/${node.slug}`}>{node.slug}</Link>{' '}
          <Metadata metadata={node.trackEventsByStationId.nodes[0]?.trackByTrackId?.metadata} />
        </div>
      ))}

      <br/>

      <CreateStation />
    </div>
  )
}

export default StationList

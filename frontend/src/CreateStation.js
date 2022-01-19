import { useState } from 'react';
import { gql, useMutation} from '@apollo/client';

const CREATE_STATION = gql`
  mutation CreateStation($slug: String!) {
    createStation(input: { slug: $slug }) {
      id
    }
  }
`

export default function CreateStation() {
  return "create me"
}

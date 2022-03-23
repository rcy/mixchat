import { NavLink, useLocation } from 'react-router-dom'

export default function SearchLink({ children, to, ...props }) {
  const { search } = useLocation();
  return (
    <NavLink to={to + search} {...props}>
      {children}
    </NavLink>
  )
}


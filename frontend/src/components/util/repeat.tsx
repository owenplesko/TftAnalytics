import { Fragment } from "react/jsx-runtime";

const Repeat: React.FC<React.PropsWithChildren<{ n: number }>> = ({
  n,
  children,
}) =>
  Array.from({ length: n }).map((_, i) => (
    <Fragment key={i}>{children}</Fragment>
  ));
export default Repeat;

import React from "react";
import { Route, Routes } from "react-router-dom";
import MainCartPage from "../Pages/cartPage/MainCartPage";
import Home from "../Pages/Home/Home";
import Login from "../Pages/Login/Login";
import Products from "../Pages/Products/Products";
import SingleProduct from "../Pages/SingleProduct/SingleProduct";
import Wishlist from "../Pages/Wishlist/Wishlist";
import Payments from "../Pages/payment/Payments";
import Checkout from "../Pages/checkout/checkout";
import { LastPage } from "../Pages/cartPage/LastPage";
import Signup from "../Pages/Signup/Signup";
import PrivateRoute from "./PrivateRoute/PrivateRoutes";
import MainLayout from "./Admin/MainLayout";
import Dashboard from "../Pages/Admin/Dashboard";
import Enquiries from "../Pages/Admin/Enquiries";
import ViewEnq from "../Pages/Admin/ViewEnq";
import Couponlist from "../Pages/Admin/Couponlist";
import AddCoupon from "../Pages/Admin/AddCoupon";
import Orders from "../Pages/Admin/Orders";
import ViewOrder from "../Pages/Admin/ViewOrder";
import Customers from "../Pages/Admin/Customers";
import Categorylist from "../Pages/Admin/Categorylist";
import Addcat from "../Pages/Admin/Addcat";
import Brandlist from "../Pages/Admin/Brandlist";
import Addbrand from "../Pages/Admin/Addbrand";
import Productlist from "../Pages/Admin/Productlist";
import Addproduct from "../Pages/Admin/Addproduct";

const AllRoutes = () => {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Home />}></Route>
        <Route
          path="/mobilesandtablets"
          element={<PrivateRoute><Products typeOfProduct="mobilesandtablets" /></PrivateRoute>}
        ></Route>
        <Route
          path="/mobilesandtablets/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="mobilesandtablets" /></PrivateRoute>}
        ></Route>
        <Route
          path="/televisions"
          element={<PrivateRoute><Products typeOfProduct="televisions" /></PrivateRoute>}
        ></Route>
        <Route
          path="/televisions/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="televisions" /></PrivateRoute>}
        ></Route>
        <Route
          path="/kitchen"
          element={<PrivateRoute><Products typeOfProduct="kitchen" /></PrivateRoute>}
        ></Route>
        <Route
          path="/kitchen/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="kitchen" /></PrivateRoute>}
        ></Route>
        <Route
          path="/headphones"
          element={<PrivateRoute><Products typeOfProduct="headphones" /></PrivateRoute>}
        ></Route>
        <Route
          path="/headphones/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="headphones" /></PrivateRoute>}
        ></Route>
        <Route
          path="/homeappliances"
          element={<PrivateRoute><Products typeOfProduct="homeappliances" /></PrivateRoute>}
        ></Route>
        <Route
          path="/homeappliances/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="homeappliances" /></PrivateRoute>}
        ></Route>
        <Route
          path="/computers"
          element={<PrivateRoute><Products typeOfProduct="computers" /></PrivateRoute>}
        ></Route>
        <Route
          path="/computers/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="computers" /></PrivateRoute>}
        ></Route>
        <Route
          path="/cameras"
          element={<PrivateRoute><Products typeOfProduct="cameras" /></PrivateRoute>}
        ></Route>
        <Route
          path="/cameras/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="cameras" /></PrivateRoute>}
        ></Route>
        <Route
          path="/personalcare"
          element={<PrivateRoute><Products typeOfProduct="personalcare" /></PrivateRoute>}
        ></Route>
        <Route
          path="/personalcare/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="personalcare" /></PrivateRoute>}
        ></Route>
        <Route
          path="/accessories"
          element={<PrivateRoute><Products typeOfProduct="accessories" /></PrivateRoute>}
        ></Route>
        <Route
          path="/accessories/:id"
          element={<PrivateRoute><SingleProduct typeOfProduct="accessories" /></PrivateRoute>}
        ></Route>
        <Route path="/cart" element={<PrivateRoute><MainCartPage /></PrivateRoute>}></Route>
        <Route path="/login" element={<Login />}></Route>
        <Route
          path="/whishlist"
          element={<PrivateRoute><Wishlist typeOfProduct={"whishlist"} /></PrivateRoute>}
        ></Route>
        <Route path="/checkout" element={<PrivateRoute><Checkout /></PrivateRoute>}></Route>
        <Route path="/payments" element={<PrivateRoute><Payments /></PrivateRoute>}></Route>
        <Route path="/lastpage" element={<PrivateRoute><LastPage /></PrivateRoute>}></Route>
        <Route path="/signup" element={<Signup />}></Route>
        <Route path="/admin" element={<MainLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="enquiries" element={<Enquiries />} />
          <Route path="enquiries/:id" element={<ViewEnq />} />
          <Route path="coupon-list" element={<Couponlist />} />
          <Route path="coupon" element={<AddCoupon />} />
          <Route path="coupon/:id" element={<AddCoupon />} />
          {/* <Route path="blog-category-list" element={<Blogcatlist />} />
          <Route path="blog-category" element={<Addblogcat />} />
          <Route path="blog-category/:id" element={<Addblogcat />} /> */}
          <Route path="orders" element={<Orders />} />
          <Route path="order/:id" element={<ViewOrder />} />
          <Route path="customers" element={<Customers />} />
          <Route path="list-category" element={<Categorylist />} />
          <Route path="category" element={<Addcat />} />
          <Route path="category/:id" element={<Addcat />} />
          <Route path="list-brand" element={<Brandlist />} />
          <Route path="brand" element={<Addbrand />} />
          <Route path="brand/:id" element={<Addbrand />} />
          <Route path="list-product" element={<Productlist />} />
          <Route path="product" element={<Addproduct />} />
        </Route>
        {/* <Route path="/order" element={<Products typeOfProduct={"order"}/>}></Route>
            <Route path="/contactus" element={<Products typeOfProduct={"contactus"}/>}></Route>
            <Route path="/profile" element={<Products typeOfProduct={"profile"}/>}></Route> */}
      </Routes>
    </div>
  );
};
export default AllRoutes;
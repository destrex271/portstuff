"use client"
// pages/index.js
import { useState, useEffect } from 'react'
import Head from 'next/head'
import dynamic from 'next/dynamic'

// Dynamically import the map component to avoid SSR issues
const MapWithNoSSR = dynamic(() => import('../components/Map'), {
  ssr: false
})

// 0 -> Small Truck
// 1 -> Meidium Sized Truck
// 2 -> Large Truck

export default function Home() {
  const [pickup, setPickup] = useState({ lat: 51.505, lng: -0.09 })
  const [dropoff, setDropoff] = useState({ lat: 51.51, lng: -0.1 })
  const [goodsDescription, setGoodsDescription] = useState('')
  const [vehicleType, setVehicleType] = useState(0)
  const [date, setDate] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    console.log({ pickup, dropoff, goodsDescription, date })
    alert('Booking submitted!')
  }

  const handleChange = (e) => {
    console.log(e.target.value)
    setVehicleType(e.target.value)
  }

  return (
    <div className="container mx-auto p-4">
      <Head>
        <title>Goods Transport Booking</title>
        <link rel="icon" href="/favicon.ico" />
        <link
          rel="stylesheet"
          href="https://unpkg.com/leaflet@1.7.1/dist/leaflet.css"
          integrity="sha512-xodZBNTC5n17Xt2atTPuE1HxjVMSvLVW9ocqUKLsCC5CXdbqCmblAshOMAS6/keqq/sMZMZ19scR4PsZChSR7A=="
          crossOrigin=""
        />
      </Head>

      <main>
        <h1 className="text-3xl font-bold mb-4">Book a Ride for Your Goods</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="pickup" className="block mb-1">Pickup Location</label>
                <input
                  type="text"
                  id="pickup"
                  value={`${pickup.lat.toFixed(4)}, ${pickup.lng.toFixed(4)}`}
                  readOnly
                  className="w-full p-2 border rounded"
                />
              </div>
              <div>
                <label htmlFor="dropoff" className="block mb-1">Dropoff Location</label>
                <input
                  type="text"
                  id="dropoff"
                  value={`${dropoff.lat.toFixed(4)}, ${dropoff.lng.toFixed(4)}`}
                  readOnly
                  className="w-full p-2 border rounded"
                />
              </div>
              <div>
                <label htmlFor="goods" className="block mb-1">Goods Description</label>
                <textarea
                  id="goods"
                  value={goodsDescription}
                  onChange={(e) => setGoodsDescription(e.target.value)}
                  required
                  className="w-full p-2 border rounded"
                ></textarea>
              </div>
              <div>
                <label htmlFor="vehicleType" className="block mb-1">Veichle Type</label>
                <select onChange={handleChange}>
                  <option value={0} key={"Small Vehicle"} />
                  <option value={1} key={"Medium Vehicle"}/>
                  <option value={2} key={"Large Vehicle"}/>
                </select>
              </div>
              <div>
                <label htmlFor="date" className="block mb-1">Pickup Date</label>
                <input
                  type="date"
                  id="date"
                  value={date}
                  onChange={(e) => setDate(e.target.value)}
                  required
                  className="w-full p-2 border rounded"
                />
              </div>
              <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
                Book Ride
              </button>
            </form>
          </div>
          <div className="h-96">
            <MapWithNoSSR pickup={pickup} setPickup={setPickup} dropoff={dropoff} setDropoff={setDropoff} />
          </div>
        </div>
      </main>
    </div>
  )
}
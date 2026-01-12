import { useEffect, useState } from 'react';
import type { InventoryItem } from '../types';

export default function InventoryDashboard() {
    const [inventory, setInventory] = useState<InventoryItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [filter, setFilter] = useState<string>('All');

    useEffect(() => {
        fetch('http://localhost:8080/api/inventory')
            .then((res) => res.json())
            .then((data) => {
                setInventory(data);
                setLoading(false);
            })
            .catch((err) => {
                console.error('Error fetching inventory:', err);
                setLoading(false);
            });
    }, []);

    // Derive unique categories
    const categories = ['All', ...new Set(inventory.map((item) => item.category))];

    // Filter items
    const filteredInventory =
        filter === 'All'
            ? inventory
            : inventory.filter((item) => item.category === filter);

    if (loading) {
        return <div className="p-8 text-center text-lg">Loading inventory...</div>;
    }

    return (
        <div className="container mx-auto p-6">
            <h1 className="text-3xl font-bold mb-6 text-gray-800">Inventory Dashboard</h1>

            {/* Filter */}
            <div className="mb-6 flex items-center gap-4">
                <label htmlFor="category-filter" className="font-medium text-gray-700">
                    Filter by Category:
                </label>
                <select
                    id="category-filter"
                    value={filter}
                    onChange={(e) => setFilter(e.target.value)}
                    className="border border-gray-300 rounded-md px-3 py-2 bg-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                    {categories.map((cat) => (
                        <option key={cat} value={cat}>
                            {cat}
                        </option>
                    ))}
                </select>
            </div>

            {/* Table */}
            <div className="overflow-x-auto shadow-md sm:rounded-lg">
                <table className="min-w-full text-sm text-left text-gray-500">
                    <thead className="text-xs text-gray-700 uppercase bg-gray-50">
                        <tr>
                            <th scope="col" className="px-6 py-3">ID</th>
                            <th scope="col" className="px-6 py-3">SKU</th>
                            <th scope="col" className="px-6 py-3">Product Name</th>
                            <th scope="col" className="px-6 py-3">Category</th>
                            <th scope="col" className="px-6 py-3">Price</th>
                            <th scope="col" className="px-6 py-3">Stock</th>
                            <th scope="col" className="px-6 py-3">Last Updated</th>
                        </tr>
                    </thead>
                    <tbody>
                        {filteredInventory.map((item) => (
                            <tr
                                key={item.id}
                                className={`border-b hover:bg-gray-50 ${item.stock === 0 ? 'bg-red-50 text-red-600 font-medium' : 'bg-white'
                                    }`}
                            >
                                <td className="px-6 py-4">{item.id}</td>
                                <td className="px-6 py-4">{item.sku}</td>
                                <td className="px-6 py-4 font-semibold">{item.product_name}</td>
                                <td className="px-6 py-4">
                                    <span className="bg-blue-100 text-blue-800 text-xs font-medium mr-2 px-2.5 py-0.5 rounded">
                                        {item.category}
                                    </span>
                                </td>
                                <td className="px-6 py-4">${item.price.toFixed(2)}</td>
                                <td className="px-6 py-4">
                                    {item.stock === 0 ? (
                                        <span className="text-red-600 font-bold">Out of Stock</span>
                                    ) : (
                                        item.stock
                                    )}
                                </td>
                                <td className="px-6 py-4">{item.last_updated}</td>
                            </tr>
                        ))}
                        {filteredInventory.length === 0 && (
                            <tr>
                                <td colSpan={7} className="px-6 py-4 text-center">
                                    No items found.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
